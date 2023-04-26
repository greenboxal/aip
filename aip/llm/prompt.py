from langchain.prompts.prompt import PromptTemplate as LlmPromptTemplate


class PromptTemplate:
    template: str
    input_variables: [str]

    def __init__(self, template: str, input_variables: [str]):
        self.template = template
        self.input_variables = input_variables

    def to_llm_template(self) -> LlmPromptTemplate:
        raise "not implemented"
