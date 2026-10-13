import { get, update } from 'idb-keyval'
import { supabase, VISIT_PHOTO_BUCKET } from './supabase'

const QUEUE_KEY = 'agritec:visit-queue'
const MAX_PHOTO_BYTES = 500 * 1024
const MAX_PHOTO_EDGE = 1280

function loadImage(file) {
    return new Promise((resolve, reject) => {
        const url = URL.createObjectURL(file)
        const img = new Image()
        img.onload = () => {
            URL.revokeObjectURL(url)
            resolve(img)
        }
        img.onerror = () => {
            URL.revokeObjectURL(url)
            reject(new Error('File foto tidak bisa dibaca.'))
        }
        img.src = url
    })
}

function canvasToBlob(canvas, quality) {
    return new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg', quality))
}

// Kompres foto ke JPEG (sisi terpanjang maks 1280px) sampai < 500 KB.
export async function compressImage(file) {
    const img = await loadImage(file)
    let scale = Math.min(1, MAX_PHOTO_EDGE / Math.max(img.width, img.height))

    for (let attempt = 0; attempt < 8; attempt++) {
        const canvas = document.createElement('canvas')
        canvas.width = Math.round(img.width * scale)
        canvas.height = Math.round(img.height * scale)
        canvas.getContext('2d').drawImage(img, 0, 0, canvas.width, canvas.height)

        for (const quality of [0.8, 0.65, 0.5]) {
            const blob = await canvasToBlob(canvas, quality)
            if (blob && blob.size <= MAX_PHOTO_BYTES) return blob
        }
        scale *= 0.75
    }
    throw new Error('Foto terlalu besar dan gagal dikompres di bawah 500 KB.')
}

// Upload foto + simpan kunjungan. Aman diulang: path foto & client_id deterministik.
async function uploadVisit(userId, visit) {
    const photoPath = `${userId}/${visit.clientId}.jpg`

    const { error: uploadError } = await supabase.storage
        .from(VISIT_PHOTO_BUCKET)
        .upload(photoPath, visit.photo, { contentType: 'image/jpeg', upsert: false })
    // Upload ulang setelah sync sebelumnya terputus: objek sudah ada, lanjutkan.
    if (uploadError && uploadError.statusCode !== '409' && !/exists/i.test(uploadError.message)) {
        throw uploadError
    }

    const { error } = await supabase.rpc('submit_visit', {
        p_client_id: visit.clientId,
        p_farmer_name: visit.farmerName,
        p_crop_type: visit.cropType,
        p_land_area: visit.landArea,
        p_latitude: visit.latitude,
        p_longitude: visit.longitude,
        p_photo_path: photoPath,
        p_notes: visit.notes || null,
        p_visit_time: visit.visitTime
    })
    if (error) throw error
}

export async function getQueue() {
    return (await get(QUEUE_KEY)) || []
}

function enqueue(visit) {
    return update(QUEUE_KEY, queue => [...(queue || []), visit])
}

function isNetworkError(err) {
    return !navigator.onLine || err instanceof TypeError || /fetch|network/i.test(err?.message || '')
}

// Kirim kunjungan; jika offline / jaringan gagal, simpan ke antrean IndexedDB.
// Mengembalikan 'sent' atau 'queued'.
export async function submitVisit(userId, form) {
    const visit = {
        clientId: crypto.randomUUID(),
        userId,
        farmerName: form.farmerName.trim(),
        cropType: form.cropType.trim(),
        landArea: Number(form.landArea),
        latitude: form.latitude,
        longitude: form.longitude,
        notes: form.notes,
        visitTime: new Date().toISOString(),
        photo: await compressImage(form.photo)
    }

    if (!navigator.onLine) {
        await enqueue(visit)
        return 'queued'
    }
    try {
        await uploadVisit(userId, visit)
        return 'sent'
    } catch (err) {
        if (!isNetworkError(err)) throw err
        await enqueue(visit)
        return 'queued'
    }
}

let syncing = null

// Kirim ulang antrean milik user yang sedang login. Mengembalikan sisa antrean.
export function syncQueue(userId) {
    if (syncing) return syncing
    syncing = (async () => {
        const sent = new Set()
        for (const visit of await getQueue()) {
            if (visit.userId !== userId || !navigator.onLine) continue
            try {
                await uploadVisit(userId, visit)
                sent.add(visit.clientId)
            } catch (err) {
                console.error('Sinkronisasi kunjungan gagal', err)
            }
        }
        // Baca ulang agar kunjungan yang masuk antrean selama sync tidak hilang.
        await update(QUEUE_KEY, queue => (queue || []).filter(v => !sent.has(v.clientId)))
        return (await getQueue()).filter(v => v.userId === userId).length
    })().finally(() => {
        syncing = null
    })
    return syncing
}
