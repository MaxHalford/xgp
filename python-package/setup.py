from codecs import open
from os import path
from setuptools import Extension
from setuptools import find_packages
from setuptools import setup

here = path.abspath(path.dirname(__file__))

with open(path.join(here, 'README.md'), encoding='utf-8') as f:
    long_description = f.read()

setup(
    name='xgp',
    version='0.0.1',
    description='A machine learning tool based on genetic programming',
    long_description=long_description,
    url='https://github.com/pypa/sampleproject',
    author='Max Halford',
    author_email='maxhalford25@gmail.com',
    license='MIT',
    classifiers=[
        # How mature is this project? Common values are
        #   3 - Alpha
        #   4 - Beta
        #   5 - Production/Stable
        'Development Status :: 3 - Alpha',
        'Intended Audience :: Developers',
        'Topic :: Software Development :: Build Tools',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
    ],
    keywords='Genetic programming',
    packages=find_packages(exclude=['tests']),
    python_requires='>=3',
    build_golang={'root': 'github.com/MaxHalford/xgp'},
    ext_modules=[Extension('xgp', ['xgp.go'])],
    setup_requires=['setuptools-golang'],
    install_requires=['numpy', 'sklearn'],
    extras_require={
        'dev': find_packages(exclude=['tests']) + ['twine', 'wheel'],
        'test': find_packages(exclude=['tests']),
    },
    entry_points={
        'console_scripts': [
            'sample=sample:main',
        ],
    },
)
