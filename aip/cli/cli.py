import click
from click_repl import register_repl

from aip.cli.index import index
from aip.cli.chat import chat
from aip.cli.generate import generate
from aip.cli.ipc import ipc

@click.group()
def root():
    pass

root.add_command(index)
root.add_command(chat)
root.add_command(generate)
root.add_command(ipc)

register_repl(root)
