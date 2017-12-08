# How it works

Symbolic regression is a supervised learning method where the shape of the model is included in the search space. The idea is to not constrain the model to have a predetermined shape and instead to include the shape in the optimization process. For example in linear regression the model is constrained to have a certain shape.

$$y = w_0x_0 + w_1x_1 + b$$

In symbolic regression the model can take many shapes. For example the following equation is valid in symbolic regression.

$$y = cos(x_0) \times log(sin(x_1)) \times \frac{3.14}{x_0} + 2.17$$

*Anything goes*. Symbolic regression is thus a very flexible approach. The obvious downside is that the search space becomes huge and cannot be explored with gradient-based optimization techniques. Even though linear regression and it's variants have a fixed structure, they are relatively simple to optimize.

A symbolic regression model is composed of three different kinds of so-called **operators**:

- **Functions** which are basic mathematical functions such as \\(cos\\) and \\(log\\),
- **Constants** which are simply floating point values such as 3.14 and 2.17,
- **Variables** which are features in a dataset such as \\(x_0\\) and \\(x_1\\).

The goal of symbolic regression is to combine these operators in an optimal way. Not only do appropriate operators have to be chosen, they also have to associated in a good way.
