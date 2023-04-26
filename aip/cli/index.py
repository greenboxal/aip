import os

import click

from aip.indexing import Index

from langchain.document_loaders import (
    DirectoryLoader,
    TextLoader,
)

from langchain.text_splitter import (
    TokenTextSplitter,
)


@click.command("create", short_help="Create index")
@click.argument("name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
def index_create(name):
    indexer = Index(index_name=name)

    print("Creating index %s" % name)
    indexer.create_index()

    print("Index %s created" % name)


@click.command("truncate", short_help="Truncate index (delete all documents)")
@click.argument("name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
@click.option("--namespace", default=None)
def index_truncate(name, namespace):
    indexer = Index(index_name=name)

    print("Truncating index %s (ns = %s)" % (name, namespace))
    indexer.truncate()

    print("Index %s truncated (ns = %s)" % (name, namespace))


@click.command("add", short_help="Indexes files into the vector store")
@click.argument("file", type=click.Path(exists=True))
@click.option("--index-name", "-i", default=lambda: os.environ["PINECONE_INDEX_NAME"], help="Index name")
@click.option("--namespace", "-n", default=None, help="Index namespace")
@click.option("--filter", "-f", default="**/*", help="Glob pattern to filter files to index")
@click.option("--relative-to", "-C", default="", help="Base directory to strip from file paths")
@click.option("--prefix", "-p", default="", help="Prefix to add to file paths")
@click.option("--dry-run", is_flag=True, default=False,
              help="Only lists files and chunk size (will still calculate embeddings)")
@click.option("--chunk-size", default=2048, help="Number of tokens per chunk generated from each document")
@click.option("--chunk-overlap", default=0,
              help="Number of tokens to overlap between adjacent chunks for each document (0 = no overlap)")
def index_add(index_name, relative_to, prefix, dry_run, filter, chunk_size, chunk_overlap, namespace, file):
    if index_name == "":
        index_name = os.environ["PINECONE_INDEX_NAME"]

    indexer = Index(index_name=index_name, namespace=namespace)

    if not os.path.exists(file):
        raise ValueError("File does not exist: %s" % file)

    if os.path.isfile(file):
        loader = TextLoader(file, encoding='utf-8')
    else:
        loader = DirectoryLoader(file, glob=filter, loader_cls=TextLoader)

    docs = loader.load()

    for doc in docs:
        if "source" in doc.metadata:
            source = doc.metadata["source"]
            source = os.path.relpath(source, relative_to)
            source = os.path.join(prefix, source)
            doc.metadata["source"] = source

        print(doc.metadata)

    text_splitter = TokenTextSplitter(chunk_size=chunk_size, chunk_overlap=chunk_overlap)
    texts = text_splitter.split_documents(docs)

    print("Importing %d documents (%d chunks)" % (len(docs), len(texts)))

    if dry_run:
        print("Dry run, not importing documents")
    else:
        indexer.import_documents(texts)
        print("Done")


@click.group("index", short_help="Index management commands")
def index():
    pass


index.add_command(index_create)
index.add_command(index_truncate)
index.add_command(index_add)
