from langchain_huggingface import HuggingFaceEmbeddings
import logging

class EmbeddingService:
    def __init__(self, model_name="sentence-transformers/all-MiniLM-L6-v2") -> None:
        self.model = HuggingFaceEmbeddings(model_name=model_name)

    def generate(self, texts: list[str]) -> list[list[float]]:
        return self.model.embed_documents(texts=texts)

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def test_generation():
    logger.info("Initializing EmbeddingService for testing...")
    service = EmbeddingService()
    
    texts = ["Teste de embedding", "IA é incrível"]
    
    logger.info(f"Generating embeddings for {len(texts)} texts...")
    vectors = service.generate(texts)
    
    logger.info(f"Successfully generated {len(vectors)} vectors.")
    
    # Validações
    dim = len(vectors[0])
    logger.info(f"First vector dimension: {dim}")
    
    assert isinstance(vectors[0][0], float), "Vector values must be of type float"
    
    logger.info("Test passed successfully: Embedding generated correctly!")

if __name__ == "__main__":
    try:
        test_generation()
    except Exception as e:
        logger.error(f"Test failed with error: {e}", exc_info=True)