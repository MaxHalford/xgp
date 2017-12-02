conda create -n mkdocs python=3.5
source activate mkdocs
pip install mkdocs mkdocs-material
mkdocs build --verbose --clean --strict
