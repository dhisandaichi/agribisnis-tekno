# ==========================================
# Batch Processing & Analytics Service (Python)
# ==========================================
# Script ini akan membaca data dari PostgreSQL
# menggunakan Pandas/Geopandas untuk di-push ke
# dashboard visualisasi atau diolah datanya (ETL).

import pandas as pd
# import geopandas as gpd
import psycopg2
import os

def check_db_connection():
    try:
        conn = psycopg2.connect(os.getenv("DATABASE_URL"))
        print("Connected to PostgreSQL successfully.")
        conn.close()
    except Exception as e:
        print(f"Error connecting to database: {e}")

if __name__ == "__main__":
    print("ETL Batch Analytics Service")
    # check_db_connection()
