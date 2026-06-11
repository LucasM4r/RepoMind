import grpc
from concurrent import futures
# Importe as classes geradas
from pb import ai_service_pb2, ai_service_pb2_grpc

# Importe suas lógicas de negócio
from internal.embeddings import EmbeddingService

class AIService(ai_service_pb2_grpc.EmbeddingServiceServicer, 
                ai_service_pb2_grpc.LLMServiceServicer):
    
    def __init__(self):
        self.embedding_logic = EmbeddingService()

    def GenerateEmbeddings(self, request, context):
        texts = request.texts

        vectors = self.embedding_logic.generate(texts=texts)
        grpc_vectors = [
            ai_service_pb2.Vector(values=v) for v in vectors
        ]

        return ai_service_pb2.EmbeddingResponse(embeddings=grpc_vectors)