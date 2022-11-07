# flask auth-app api

## Run Locally

Clone the project

```bash
  git clone https://github.com/kidboy-man/white-whale-api.git
```

Go to the project directory. Install [pyenv](https://github.com/pyenv/pyenv) 

```bash
  cd white-whale-api/auth-app
  pyenv install 3.10.8
  pyenv shell 3.10.8
  pip install virtualenv
  virtualenv .venv
  source .venv/bin/activate
  pip install -r requirements.txt
  flask run
```