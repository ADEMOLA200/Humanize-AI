import os
from dotenv import load_dotenv
from flask import Flask, request, jsonify
from transformers import T5ForConditionalGeneration, T5Tokenizer
from huggingface_hub import login

load_dotenv()

HF_TOKEN = os.getenv("HF_TOKEN")

if not HF_TOKEN:
    raise ValueError("Hugging Face token not found. Set HF_TOKEN as an environment variable or in a .env file.")

login(HF_TOKEN)

app = Flask(__name__)

# If you later want to upgrade to Flan-T5-XXL, you could uncomment the lines below:
# model_name = "google/flan-t5-xxl"
# from transformers import AutoTokenizer, AutoModelForSeq2SeqLM
# tokenizer = AutoTokenizer.from_pretrained(model_name)
# model = AutoModelForSeq2SeqLM.from_pretrained(model_name)

# Load model and tokenizer with authentication, you can uncomment any one
# and use them, but "Vamsi/T5_Paraphrase_Paws" is simply preferable so yup!
# model_name = "t5-small"
# model_name = "t5-base"
model_name = "Vamsi/T5_Paraphrase_Paws"
model = T5ForConditionalGeneration.from_pretrained(model_name, token=HF_TOKEN)
tokenizer = T5Tokenizer.from_pretrained(model_name, token=HF_TOKEN)

@app.route('/paraphrase', methods=['POST'])
def paraphrase():
    data = request.get_json()
    text = data.get('text', '')
    if not text:
        return jsonify({'error': 'Text is required.'}), 400

    instruction = "Humanize this AI-generated text into a casual, conversational style: "
    prompt = instruction + text

    inputs = tokenizer(prompt, return_tensors="pt", max_length=768, truncation=True)
    
    outputs = model.generate(
        inputs.input_ids,
        max_length=768,
        min_length=200,  # ensures at least 200 tokens are generated
        do_sample=True,
        num_beams=5,
        num_return_sequences=1,
        temperature=0.7,  # Lower temperature for more controlled output
        top_p=0.9,
        top_k=50
    )
    
    generated_text = tokenizer.decode(outputs[0], skip_special_tokens=True)
    return jsonify({'paraphrased': generated_text, 'success': True})

if __name__ == '__main__':
    app.run(port=5001)
