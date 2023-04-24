# Crazy AI Stuff™

## Setup

```sh
# Create python venv
python3 -m venv "venv"

# Activate venv (once per shell)
source ./venv/bin/activate

# Install requirements
pip install -r requirements.txt

# Copy over the .env template
cp .env-template .env
```

Finally, make sure you populate your `.env` file with the correct API keys.

## Running

```sh
# Activate venv (once per shell)
source ./venv/bin/activate

# Load env variables (once per shell)
source .env

# Run main CLI tool
python -m aip --help
```
