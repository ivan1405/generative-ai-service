# Generative AI service

This service allows users to easily integrate with different Generative AI providers. For now, the following providers are supported:
- Chat GPT
- AWS Bedrock

It exposes an API with the following endpoints:

## Get AI capabilities

Not every provider enables the same GenAI functionalities. Some of them may be able to create text-to-speech or speech-to-text, but some other may not.
There is an endpoint that exposes all the capabilities integrated in the service of each of the providers

### Request

```GET /ai-capabilities```

### Response
```json
{
    "aws-bedrock": [
        "Completions",
        "Image generation"
    ],
    "chat-gpt": [
        "Completions",
        "Image generation"
    ]
}
```

## Create chat completion

Allows users to create chat completions by providing a prompt as input paramter. 

### Request

Payload

```json
{
    "prompt": "What's the purpose of life?",
    "model": "gpt-5",    //optional
    "temperature": 0.1,  //optional
    "max_tokens": 200,   //optional
    "top_p": 0.9         //optional
}
```

If optional parameters are not informed, the default ones indicated by each provider will be set

Header 'Provider' can also be set to indicate where to redirect our request, possible values are:
- Chat GPT -> `Provider: 'chat-gpt'`
- AWS Bedrock -> `Provider: 'aws-bedrock'`

If no header is set, default provider is Chat GPT

### Response
```json
{
    "response": "As an AI, I don't have personal experiences or beliefs. However, I can tell you that the purpose of life is a philosophical and existential question that has been debated by scholars, theologians, and thinkers throughout history. Some believe the purpose of life is to seek happiness, knowledge, or spiritual enlightenment. Others believe it's to contribute to the betterment of humanity, to love and be loved, or to express oneself creatively. Ultimately, the purpose of life may be a deeply personal and subjective concept that differs from person to person.",
    "provider": "chat-gpt"
}
```

## Generate images

Image generation is one of the most famous features of Generative AI providers. This has been also integrated in the service.

### Request

Payload

```json
{
    "prompt": "A picture of a happy dog playing in the beach",
}
```

### Response
The response will return the image encoded in Base64
```json
{
    "image": "iVBORw0KGgoAAAANSUhEUgAAAgAAAAIACAIAAAB75+f8HBAX...tA7hv6xkuJwAAAABJRU5ErkJggg==",
    "provider": "chat-gpt"
}
```


## Configuration

In order for the service to work, provider credentials need to be configured. These can be set easily as environment variables in a `.env` file placed in the root folder.

```
CHAT_GPT_KEY="sk-dsd<szxdf"
AWS_ACCESS_KEY_ID="DSFDSAFDASFHSDSDFG"
AWS_SECRET_ACCESS_KEY ="dsdfSD43sadF4asdfWEdsfafda$/dcf"
```
