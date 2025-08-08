import json
import time
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed

def lambda_handler(event, context):
    try:
        start_time = time.time()
        
        def count_worker(thread_id):
            for i in range(1, 101):
                print(f"Thread {thread_id}: {i}")
        
        with ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(count_worker, i) for i in range(10)]
            
            # Wait for all threads to complete
            for future in as_completed(futures):
                future.result()
        
        end_time = time.time()
        duration = (end_time - start_time) * 1000  # Convert to milliseconds
        
        response = {
            'statusCode': 200,
            'body': json.dumps({
                'language': 'Python',
                'threads': 10,
                'count_per_thread': 100,
                'duration_ms': duration,
                'status': 'completed'
            })
        }
        
    except Exception as e:
        response = {
            'statusCode': 500,
            'body': json.dumps({
                'error': str(e),
                'status': 'failed'
            })
        }
    
    return response