from langchain_huggingface import HuggingFaceEmbeddings

class EmbeddingService:
    def __init__(self, model_name="sentence-transformers/all-MiniLM-L6-v2") -> None:
        self.model = HuggingFaceEmbeddings(model_name=model_name)

    def generate(self, texts: list[str]) -> list[list[float]]:
        return self.model.embed_documents(texts=texts)