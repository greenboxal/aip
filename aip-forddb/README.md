# FordDB

FordDB is a document-based database designed with git-like semantics for branches and commit (change) histories. It's built on go-ipld-prime, which makes it a mutable, concurrent database that supports dynamic linking to other objects and optimistic locking. 

## Features
- Git-like semantics like branches and commit histories.
- Supports Get, List, Update, Delete operations.
- Support for both static (IPLD Link) and dynamic (by object ID) references to other objects.
- Optimistic locking for Update and Delete operations.
- Ability to retrieve a specific version of a document either by version or by IPLD link.
- Support for arbitrary filters and sorting based on the document IPLD schema in List operations.

## Directory Structure

Here is a brief overview of our project directory:

- `ipld/`: Contains modules for managing IPLD data structures.
  - `schemas/`: Manages IPLD schemas.
  - `linking/`: Handles IPLD Links.
  - `codec/`: Handles different codecs used in IPLD.
  - `store/`: Handles the storage and retrieval of IPLD data.

- `documents/`: Contains modules for managing document operations.
  - `models/`: Manages document models.
  - `indices/`: Manages indices on documents.
  - `operations/`: Handles document CRUD operations.
  - `utils/`: Utilities for handling document operations.

- `concurrency/`: Contains modules for managing concurrent operations.

- `versioning/`: Contains modules for managing version control aspects.

- `tests/`: Contains test files and test data.

## Usage

Detailed usage guide will be updated soon.

## License

FordDB is an open-source project under [license name].

## Contributions

We welcome contributions from the community. Please refer to our [CONTRIBUTING.md](./CONTRIBUTING.md) for more details.

## Contact

For any queries, feel free to reach out at [contact email].