import os
from dotenv import load_dotenv
from flask import Flask, request, jsonify
from transformers import T5ForConditionalGeneration, T5Tokenizer
from huggingface_hub import login
import torch

load_dotenv()

HF_TOKEN = os.getenv("HF_TOKEN")

if not HF_TOKEN:
    raise ValueError("Hugging Face token not found. Set HF_TOKEN as an environment variable or in a .env file.")

login(HF_TOKEN)

app = Flask(__name__)

# Load model and tokenizer with authentication, you can uncomment any one
# and use them, but "Vamsi/T5_Paraphrase_Paws" is simply preferable so yup!
# model_name = "t5-small"
# model_name = "google/flan-t5-large"
# model_name = "t5-base"
model_name = "Vamsi/T5_Paraphrase_Paws"


model = T5ForConditionalGeneration.from_pretrained(model_name, token=HF_TOKEN)
tokenizer = T5Tokenizer.from_pretrained(model_name, token=HF_TOKEN)

@app.route('/paraphrase', methods=['POST'])
def paraphrase():
    text = request.json['text']
    inputs = tokenizer.encode("paraphrase: " + text, return_tensors="pt", max_length=768, truncation=True)
    outputs = model.generate(
        inputs,
        max_length=768,
        min_length=200,  # this will ensures at least 200 tokens are generated
        do_sample=True,
        num_beams=5,
        num_return_sequences=1,
        temperature=1.5,
        top_p=0.9
    )
    return jsonify({'paraphrased': tokenizer.decode(outputs[0], skip_special_tokens=True)})

if __name__ == '__main__':
    app.run(port=5001)
