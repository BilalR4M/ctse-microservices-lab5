import os
import time
import requests
import psycopg2
from psycopg2.extras import RealDictCursor
from flask import Flask, request, jsonify

app = Flask(__name__)

# Wait and Connect DB
def get_db_connection():
    db_host = os.environ.get('DB_HOST', 'postgres')
    db_port = os.environ.get('DB_PORT', '5432')
    db_user = os.environ.get('DB_USER', 'postgres')
    db_pass = os.environ.get('DB_PASS', 'password')
    db_name = os.environ.get('DB_NAME', 'microservices_db')

    for i in range(5):
        try:
            conn = psycopg2.connect(
                host=db_host,
                port=db_port,
                user=db_user,
                password=db_pass,
                dbname=db_name
            )
            return conn
        except Exception as e:
            print(f"Retrying connection to DB... ({i+1}/5)")
            time.sleep(5)
    raise Exception("Could not connect to database")

# Test connection on startup
try:
    conn = get_db_connection()
    conn.close()
    print("Connected to PostgreSQL DB.")
except:
    pass

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "UP", "service": "payment-service"}), 200

@app.route('/payments', methods=['GET'])
def get_payments():
    try:
        conn = get_db_connection()
        cur = conn.cursor(cursor_factory=RealDictCursor)
        cur.execute("SELECT id, order_id, amount, method, status, created_at FROM payment_schema.payments ORDER BY id DESC")
        payments = cur.fetchall()
        cur.close()
        conn.close()
        
        for p in payments:
            if 'created_at' in p and p['created_at']:
                p['created_at'] = p['created_at'].isoformat()
                
        return jsonify(payments), 200
    except Exception as e:
        print(e)
        return jsonify({"error": "Internal Server Error"}), 500

@app.route('/payments/<int:payment_id>', methods=['GET'])
def get_payment_by_id(payment_id):
    try:
        conn = get_db_connection()
        cur = conn.cursor(cursor_factory=RealDictCursor)
        cur.execute("SELECT id, order_id, amount, method, status, created_at FROM payment_schema.payments WHERE id = %s", (payment_id,))
        payment = cur.fetchone()
        cur.close()
        conn.close()
        
        if not payment:
            return jsonify({"error": "Payment not found"}), 404
            
        if 'created_at' in payment and payment['created_at']:
            payment['created_at'] = payment['created_at'].isoformat()
            
        return jsonify(payment), 200
    except Exception as e:
        print(e)
        return jsonify({"error": "Internal Server Error"}), 500

@app.route('/payments/process', methods=['POST'])
def process_payment():
    data = request.json
    if not data or 'order_id' not in data or 'amount' not in data:
        return jsonify({"error": "Invalid request. Need order_id and amount"}), 400
        
    order_id = data['order_id']
    amount = data['amount']
    method = data.get('method', 'CREDIT_CARD')
    
    order_service_url = os.environ.get('ORDER_SERVICE_URL', 'http://order-service:8082')
    
    # 1. Validate order exists
    # 2. Update order status to PAID
    try:
        order_resp = requests.get(f"{order_service_url}/orders/{order_id}")
        if order_resp.status_code != 200:
            return jsonify({"error": "Order ID is invalid or does not exist"}), 400
            
        update_resp = requests.put(f"{order_service_url}/orders/{order_id}/status", json={"status": "PAID"})
        if update_resp.status_code != 200:
            return jsonify({"error": "Failed to update order status"}), 500
            
    except Exception as e:
        print(e)
        return jsonify({"error": f"Failed to communicate with order service: {str(e)}"}), 500

    # 3. Save payment
    try:
        conn = get_db_connection()
        cur = conn.cursor(cursor_factory=RealDictCursor)
        cur.execute(
            "INSERT INTO payment_schema.payments (order_id, amount, method, status) VALUES (%s, %s, %s, 'SUCCESS') RETURNING *",
            (order_id, amount, method)
        )
        payment = cur.fetchone()
        conn.commit()
        cur.close()
        conn.close()
        
        if payment and 'created_at' in payment and payment['created_at']:
            payment['created_at'] = payment['created_at'].isoformat()
            
        return jsonify(payment), 201
    except Exception as e:
        print(e)
        return jsonify({"error": f"Failed to save payment: {str(e)}"}), 500

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8083))
    app.run(host='0.0.0.0', port=port)
