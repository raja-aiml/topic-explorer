# 🧭 AI Explorer  

## Overview  
**AI Explorer** is a command-line tool that generates structured prompts from YAML templates and interacts with Large Language Models (LLMs) like OpenAI's GPT and Ollama-based models. It helps users create detailed prompts for AI interactions, retrieve model responses, and save them for further analysis.  

---

## 🚀 Features  
✅ **Generate AI Prompts**: Convert structured YAML templates into well-formatted prompts  
✅ **Interact with LLMs**: Communicate with OpenAI GPT models or Ollama-based models  
✅ **Multi-Provider Support**: Works with multiple LLM providers like OpenAI and Ollama  
✅ **Custom Configuration**: Define model parameters, temperature, and timeout settings  
✅ **Save Responses**: Store generated prompts and responses for later use  
✅ **Modular Design**: Easily extendable structure with separate `cmd`, `llm`, `paths`, and `config` modules  

---

## 📦 Installation  

### Prerequisites  
Ensure you have the following installed:  
- **Go 1.24+**  
- **Homebrew (macOS users)** (for package management)  
- **Python (if using Pyenv for additional dependencies)**  

### Clone the Repository  
```sh  
git clone https://github.com/yourusername/topic-explorer.git  
cd ai-explorer  
```

### Install Dependencies  
```sh  
go mod tidy  
```

### Build the Project  
```sh  
go build -o ai-explorer  
```

---

## 🎯 Usage  

### Generate a Prompt  
```sh  
./ai-explorer prompt --topic git --template resources/templates/topic.yaml --config resources/configs/git.yaml --output resources/output/git/prompt.txt  
```

### Interact with LLM  
```sh  
./ai-explorer llm --provider openai --model gpt-4 --prompt resources/output/git/prompt.txt --temperature 0.8  
```

### Generate & Retrieve a Response in One Step  
```sh  
./ai-explorer chat --topic git --provider openai --model gpt-4o  
```

---

## ⚙️ Configuration  

### Modify Default Configurations  
Default configurations are stored in `resources/default/config.yaml`. You can customize:  
```yaml  
provider: "ollama"  

model:  
  name: "phi4"  
  temperature: 0.8  
  streaming: true  

client:  
  timeout: "2m"  
  verbose_logging: true  
```

### Add New Topics  
To create a new topic, add a new YAML file inside `resources/configs/`:  
```yaml  
audience: "Developers"  
learning_stage: "Intermediate"  
topic: "Docker"  
context: "Software Engineering"  \analogies: "Shipping containers and package management"  
```

---

## 🔌 Supported LLM Providers  

| Provider | Models Supported |  
|----------|-----------------|  
| **OpenAI** | gpt-4o, gpt-3.5-turbo |  
| **Ollama** | phi4, mistral, llama3 |  

### OpenAI Usage  
```sh  
./ai-explorer llm --provider openai --model gpt-4o  
```

### Ollama (Local Inference)  
```sh  
./ai-explorer llm --provider ollama --model phi4  
