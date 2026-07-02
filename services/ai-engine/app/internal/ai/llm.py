import os

import torch
from transformers import pipeline
from langchain_huggingface import HuggingFacePipeline


class LLMService:
    def __init__(self, model_name: str | None = None) -> None:
        model_name = model_name or os.getenv("LLM_MODEL_NAME", "TinyLlama/TinyLlama-1.1B-Chat-v1.0")
        llm_device = os.getenv("LLM_DEVICE", "auto").lower()
        if llm_device == "cpu":
            pipeline_device = -1
        elif llm_device == "cuda":
            pipeline_device = 0
        else:
            pipeline_device = 0 if torch.cuda.is_available() else -1

        pipe = pipeline(
            "text-generation",
            model=model_name,
            device=pipeline_device,
            max_new_tokens=256,
            temperature=0.1,
            repetition_penalty=1.15
        )

        self.llm = HuggingFacePipeline(pipeline=pipe)

    def generate(self, prompt: str) -> str:
        return self.llm.invoke(prompt)