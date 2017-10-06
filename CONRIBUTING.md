# Contributing

## To do

- Range selection
- Caching
- Adaboost regressor
- Add boosting to command-line
- http://www.genetic-programming.com/gpquadraticexample.html
- http://www.genetic-programming.com/gpflowchart.html
- [Classification Strategies for Image Classification in Genetic Programming](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf)
- [DRS](http://goanna.cs.rmit.edu.au/~vc/papers/loveard-phd.pdf)
- [Sampling data for fitness evaluation](http://eplex.cs.ucf.edu/papers/morse_gecco16.pdf)
- http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf
- http://cswww.essex.ac.uk/staff/poli/gp-field-guide/42StepbyStepSampleRun.html
- Consider parsimony for generalization
- Use a cache that empties every k generations to store results of sub-branches to avoid unnecessary evaluations

## Documentation

The documentation is rendered with [MkDocs](http://www.mkdocs.org/). Run the following commands to run the documentation locally.

```sh
pip install mkdocs mkdocs-material pygments
cd /path/to/xgp
mkdocs serve
```
