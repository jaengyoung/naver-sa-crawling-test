package example;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.Map;
import java.util.HashMap;

public class LambdaHandler implements RequestHandler<Map<String, Object>, Map<String, Object>> {
    
    @Override
    public Map<String, Object> handleRequest(Map<String, Object> input, Context context) {
        Map<String, Object> response = new HashMap<>();
        
        try {
            long startTime = System.currentTimeMillis();
            
            ExecutorService executor = Executors.newFixedThreadPool(10);
            CountDownLatch latch = new CountDownLatch(10);
            
            for (int threadId = 0; threadId < 10; threadId++) {
                final int tid = threadId;
                executor.submit(() -> {
                    try {
                        for (int i = 1; i <= 100; i++) {
                            System.out.println("Thread " + tid + ": " + i);
                        }
                    } catch (Exception e) {
                        Thread.currentThread().interrupt();
                    } finally {
                        latch.countDown();
                    }
                });
            }
            
            latch.await(30, TimeUnit.SECONDS);
            executor.shutdown();
            
            long endTime = System.currentTimeMillis();
            long duration = endTime - startTime;
            
            response.put("language", "Java");
            response.put("threads", 10);
            response.put("count_per_thread", 100);
            response.put("duration_ms", duration);
            response.put("status", "completed");
            
        } catch (Exception e) {
            response.put("error", e.getMessage());
            response.put("status", "failed");
        }
        
        return response;
    }
}