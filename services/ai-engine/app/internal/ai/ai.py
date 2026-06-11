import grpc
from concurrent import futures
from pb import ai_service_pb2, ai_service_pb2_grpc

from internal.ai.embeddings import EmbeddingService
from internal.ai.llm import LLMService

class AIService(ai_service_pb2_grpc.EmbeddingServiceServicer, 
                ai_service_pb2_grpc.LLMServiceServicer):
    
    def __init__(self):
        self.embedding_logic = EmbeddingService()
        self.llm_logic = LLMService()

    def GenerateEmbeddings(self, request, context):
        texts = request.texts

        vectors = self.embedding_logic.generate(texts=texts)
        grpc_vectors = [
            ai_service_pb2.Vector(values=v) for v in vectors
        ]

        return ai_service_pb2.EmbeddingResponse(embeddings=grpc_vectors)
    
    def GenerateCompletion(self, request, context):
        if len(request.history) == 0:
            return ai_service_pb2.CompletionResponse(text="Erro: Nenhuma mensagem enviada.")
            
        user_prompt = request.history[-1].content
        
        generated_text = self.llm_logic.generate(prompt=user_prompt)
        return ai_service_pb2.CompletionResponse(text=generated_text)