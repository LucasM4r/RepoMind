from transformers import pipeline
from langchain_huggingface import HuggingFacePipeline
class LLMService:
    def __init__(self, model_name="TinyLlama/TinyLlama-1.1B-Chat-v1.0") -> None:
        pipe = pipeline(
            "text-generation", 
            model=model_name, 
            max_new_tokens=256, 
            temperature=0.1,
            repetition_penalty=1.15
            )

        self.llm = HuggingFacePipeline(pipeline=pipe)
    
    def generate(self, prompt: str) -> str:
        return self.llm.invoke(prompt)