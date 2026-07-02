import os

import torch
from langchain_huggingface import HuggingFaceEmbeddings


class EmbeddingService:
    def __init__(self, model_name: str | None = None) -> None:
        model_name = model_name or os.getenv("EMBEDDING_MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2")
        embedding_device = os.getenv("EMBEDDING_DEVICE", "auto").lower()
        if embedding_device == "cpu":
            device = "cpu"
        elif embedding_device == "cuda":
            device = "cuda"
        else:
            device = "cuda" if torch.cuda.is_available() else "cpu"

        self.model = HuggingFaceEmbeddings(
            model_name=model_name,
            model_kwargs={"device": device},
        )

    def generate(self, texts: list[str]) -> list[list[float]]:
        return self.model.embed_documents(texts=texts)