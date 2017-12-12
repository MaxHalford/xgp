# How it works

Symbolic regression is a supervised learning method where the shape of the model is included in the search space. The idea is to not constrain the model to have a predetermined shape and instead to include the shape in the optimization process. For example in linear regression the model is constrained to have a certain shape.

$$y = w_0x_0 + w_1x_1 + b$$

In symbolic regression the model can take many shapes. For example the following equation is valid in symbolic regression.

$$y = cos(x_0) \times log(sin(x_1)) \times \frac{3.14}{x_0} + 2.17$$

*Anything goes*. Symbolic regression is thus a very flexible approach. The obvious downside is that the search space becomes huge and unexplorable by standard gradient-based methods. Even though linear regression and it's variants have a fixed structure, they are relatively simple to optimize.

A symbolic regression model is composed of three different kinds of so-called **operators**:

- **Functions** which are basic mathematical functions such as \\(cos\\) and \\(log\\),
- **Constants** which are simply floating point values such as 3.14 and 2.17,
- **Variables** which are features in a dataset such as \\(x_0\\) and \\(x_1\\).

Each operator, regardless of it's type, has an **arity** which determines how many operands it takes. For example the multiplication operator has an arity of 2. In general all functions have an arity of 1 or more, whilst constants and variables have an arity of 0. Arity is needed internally to determine if an equation is legal or not. For example a function with an arity of 2 should have two operands. The idea is that an equation can be represented as a tree with each node being an operator and each branch an operand. For example the following tree represents the equation presented above.

![example1](img/example1.png)

!!! info
    This diagram was generated with the [CLI's todot command](cli#todot).

The goal of symbolic regression is to find an optimal combination of operators. Not only do appropriate operators have to be chosen, they also have to associated in a good way. The search space is very complex and cannot be explored with gradient-based optimization techniques. Indeed symbolic regression relies on **evolutionary algorithms** such as **genetic algorithms** to perform the optimization phase.

!!! info
    Symbolic regression is part of the larger family of **genetic programming**. Although they are very close and related, genetic programming and genetic algorithms are different things.

In genetic programming combinations of operators are usually referred to as **programs**. Like other genetic programming methods, the idea with symbolic regression is to **evolve** said programs, hopefully by making them at solving the machine learning task as hand.

TODO

