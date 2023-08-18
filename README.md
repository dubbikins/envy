# Envy

This library is used to unmarshal values from environment variables into a struct type utilizing struct tags.

## Tag syntax

To unmarshal an environment variable into a struct field, utilize the "env" struct tag with the following context free grammar.




$$
\begin{align}
&S \to \text{env:"}E\text{"} | \text{env:"}E;A\text{"}\\
&E \to \text{a} | \text{b} | \text{c} |...| \text{Y} | \text{Z} | 0 | 1 |...| 9 | \\_\\
&A \to \text{''} | ;Optional | ;Optional; | ;Optional;Optional\\
&Options \to options=[OptionsValue]\\
&Default \to \text{default=""}\\
&Optional \to \text{options=}[ OptionsValue ] |  | \text{required}\\
&OptionsValue \to \text{""}\\
\end{align}
$$


Note: 
- `"` is a literal character in the grammar while `'` around a string is meant to represent a literal.