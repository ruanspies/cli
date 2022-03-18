# setup pip configurations for uploads
mkdir $HOME/.config/pip
echo "[global]
index-url = https://europe-west1-python.pkg.dev/alis-org-777777/protobuf-python/simple/" > $HOME/.config/pip/pip.conf

echo "[distutils]
index-servers =
    protobuf-python

[protobuf-python]
repository: https://europe-west1-python.pkg.dev/alis-org-777777/protobuf-python/" > ~/.pypirc

# remove previous package version from local directory
rm -rf dist
rm -rf alis_exchange_protobuf.egg-info

# generate latest package
python3 setup.py sdist

## upload package to artefact registry repo
# install keyring library which provides applications with a way to access keyring backends,
# meaning operating system and third-party credential stores, and install the Artifact Registry backend
pip3 install keyring
pip3 install keyrings.google-artifactregistry-auth
pip3 install twine

twine upload --repository-url https://europe-west1-python.pkg.dev/alis-org-777777/protobuf-python/ dist/*

## remove pip config's
rm -rf some_dir $HOME/.config/pip
rm -rf some_dir $HOME/.pypirc