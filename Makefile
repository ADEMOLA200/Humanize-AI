GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
dev:
	cd "/Users/ademola/Documents/project/Undetectable AI" && air

uvicorn:
	uvicorn nlp_server:app --host 0.0.0.0 --port 8000 --reload

uvicorn-reload:
	uvicorn nlp_server:app --reload

# DP SK
start-server:
	python DSk/models/paraphrase_server.py

# Install dependencies
requirements:
	pip install -r requirements.txt

# Start the server
# start-server:
# 	python models/openai_t5_server.py
