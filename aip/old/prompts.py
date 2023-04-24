from langchain.prompts import (
    ChatPromptTemplate,
    PromptTemplate,
    SystemMessagePromptTemplate,
    MessagesPlaceholder,
    HumanMessagePromptTemplate,
)
from langchain.schema import (
    AIMessage,
    HumanMessage,
    SystemMessage
)


compress_prompt = PromptTemplate(
    input_variables=["description", "data"],
    template="""Compress the following {description} in a way that is lossless but results in the minimum number of tokens 
    which could be fed into an LLM like yourself as-is and produce the same output. Feel free to use multiple 
    languages, symbols, other up-front priming to lay down rules. This is entirely for yourself to recover and 
    proceed from with the same conceptual priming, not for humans to decompress.\n{data}"""
)

decompress_prompt = PromptTemplate(
    input_variables=["description", "data"],
    template="""Decompress the following tokens as {description}:\n{data}"""
)
