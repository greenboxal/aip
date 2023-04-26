import os

import click
import readline

from aip.indexing import Index
from aip.llm import Agent


@click.command("generate", short_help="Generate code")
@click.argument("recipe", type=click.Path(exists=True))
@click.option("--index-name", default=lambda: os.environ["PINECONE_INDEX_NAME"], help="Index name")
@click.option("--namespace", default=None, help="Index namespace")
def generate(index_name, namespace):
    indexer = Index(index_name=index_name, namespace=namespace)
    retriever = indexer.as_retriever()
    agent = Agent(retriever)

    readline.parse_and_bind('set editing-mode emacs')

    while True:
        line = input(">>> ")

        if 'Exit' == line.rstrip():
            break

        result = agent.codex_chain.predict(input=line)

        print(result)

    print("Done")
