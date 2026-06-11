import grpc
from concurrent import futures
from pb import ai_service_pb2, ai_service_pb2_grpc

from internal.ai.ai import AIService

def serve():
    try:
        print("[INFO] Initializing gRPC server pool...")
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

        ai_service = AIService()

        print("[INFO] Registering EmbeddingService...")
        ai_service_pb2_grpc.add_EmbeddingServiceServicer_to_server(ai_service, server)

        print("[INFO] Registering LLMService...")
        ai_service_pb2_grpc.add_LLMServiceServicer_to_server(ai_service, server)

        server.add_insecure_port('[::]:50051')
        
        print("[INFO] gRPC Server starting on port 50051")
        server.start()
        
        print("[INFO] Server is running. Press Ctrl+C to stop.")
        server.wait_for_termination()
        
    except Exception as e:
        print(f"[ERROR] Failed to start server: {e}")
        import traceback
        traceback.print_exc()