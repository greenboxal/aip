import click

from aip.cli.index import index
from aip.cli.chat import chat

@click.group()
def root():
    pass

root.add_command(index)
root.add_command(chat)
