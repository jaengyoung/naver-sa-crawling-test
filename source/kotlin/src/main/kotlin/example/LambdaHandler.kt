package example

import com.amazonaws.services.lambda.runtime.Context
import com.amazonaws.services.lambda.runtime.RequestHandler
import kotlinx.coroutines.*

class LambdaHandler : RequestHandler<Map<String, Any>, Map<String, Any>> {
    
    override fun handleRequest(input: Map<String, Any>, context: Context): Map<String, Any> {
        return runBlocking {
            val response = mutableMapOf<String, Any>()
            
            try {
                val startTime = System.currentTimeMillis()
                
                val jobs = (0 until 10).map { coroutineId ->
                    async(Dispatchers.Default) {
                        for (i in 1..100) {
                            println("Coroutine $coroutineId: $i")
                        }
                    }
                }
                
                jobs.awaitAll()
                
                val endTime = System.currentTimeMillis()
                val duration = endTime - startTime
                
                response["language"] = "Kotlin"
                response["coroutines"] = 10
                response["count_per_coroutine"] = 100
                response["duration_ms"] = duration
                response["status"] = "completed"
                
            } catch (e: Exception) {
                response["error"] = e.message ?: "Unknown error"
                response["status"] = "failed"
            }
            
            response
        }
    }
}