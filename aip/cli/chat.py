import os
import sys

import click
import readline

from aip.indexing import Indexer
from aip.llm import Agent


@click.command("chat")
@click.option("--index-name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
@click.option("--namespace", default=None)
def chat(index_name, namespace):
    indexer = Indexer(index_name=index_name, namespace=namespace)
    retriever = indexer.as_retriever()
    agent = Agent(retriever)

    readline.parse_and_bind('set editing-mode emacs')

    while True:
        line = input(">>> ")

        if 'Exit' == line.rstrip():
            break

        try:
            result = agent.codex_chain.predict(input=line)
            print(result)
        except Exception as err:
            print("### EXCEPTION ###")
            print(err)

    print("Done")
