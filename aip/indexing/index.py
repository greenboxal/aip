import os

import pinecone
from langchain.vectorstores import Milvus
from langchain.embeddings.openai import OpenAIEmbeddings

class Index:
    def __init__(
            self,
            index_name,
            namespace=None,
    ):
        self.embeddings = OpenAIEmbeddings()
        self.index_name = index_name
        self.namespace = namespace
        self.vs = None

    def create_index(self):
        pinecone.create_index(
            name=self.index_name,
            dimension=1536,  # dimensionality of dense model
            metric="dotproduct",  # sparse values supported only for dotproduct
            pod_type="s1",
            metadata_config={"indexed": []}  # see explaination above
        )

    def get_vector_store(self):
        if self.vs is None:
            self.vs = Milvus(
                collection_name="memories",
                connection_args={
                    "host": os.environ["MILVUS_HOST"],
                    "port": os.environ["MILVUS_PORT"],
                    "user": os.environ["MILVUS_USERNAME"],
                    "password": os.environ["MILVUS_PASSWORD"],
                    "secure": True,
                },
                embedding_function=self.embeddings,
            )

        return self.vs

    def as_retriever(self, **kwargs):
        return self.get_vector_store().as_retriever(**kwargs)

    def import_documents(self, documents):
        Milvus.from_documents(
            collection_name="memories",
            connection_args={
                "host": os.environ["MILVUS_HOST"],
                "port": os.environ["MILVUS_PORT"],
                "user": os.environ["MILVUS_USER"],
                "password": os.environ["MILVUS_PASSWORD"],
                "secure": True,
            },
            embedding_function=self.embeddings,
            documents=documents,
        )

    def truncate(self):
        index = pinecone.Index(self.index_name)

        index.delete(
            delete_all = True,
            namespace = self.namespace,
        )


if os.environ["PINECONE_API_KEY"] is not None:
    pinecone.init(api_key=os.environ["PINECONE_API_KEY"], environment=os.environ["PINECONE_ENV"])
