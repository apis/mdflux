# KaTeX Scientific Syntax Examples

A comprehensive test document for KaTeX mathematical rendering.

---

## 1. Inline Math

The quadratic formula is $x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$ and Einstein's famous equation is $E = mc^2$.

The area of a circle is $A = \pi r^2$ where $r$ is the radius.

---

## 2. Display Math (Block Equations)

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$

$$
\sum_{n=1}^{\infty} \frac{1}{n^2} = \frac{\pi^2}{6}
$$

---

## 3. Fractions and Binomials

$$
\frac{n!}{k!(n-k)!} = \binom{n}{k}
$$

$$
\cfrac{1}{1 + \cfrac{1}{1 + \cfrac{1}{1 + x}}}
$$

---

## 4. Greek Letters

| Lowercase | Symbol | Uppercase | Symbol |
|-----------|--------|-----------|--------|
| alpha | $\alpha$ | Alpha | $A$ |
| beta | $\beta$ | Beta | $B$ |
| gamma | $\gamma$ | Gamma | $\Gamma$ |
| delta | $\delta$ | Delta | $\Delta$ |
| epsilon | $\epsilon$ | — | — |
| theta | $\theta$ | Theta | $\Theta$ |
| lambda | $\lambda$ | Lambda | $\Lambda$ |
| sigma | $\sigma$ | Sigma | $\Sigma$ |
| omega | $\omega$ | Omega | $\Omega$ |

---

## 5. Subscripts and Superscripts

$$
x_1, x_2, \ldots, x_n \quad \text{and} \quad a^{2^n}
$$

$$
\tensor*[^{14}_{6}]{\text{C}}{}
$$

$$
x_{i,j}^{2} + y_{i,j}^{2} = z_{i,j}^{2}
$$

---

## 6. Integrals

$$
\int_0^1 x^2 \, dx = \frac{1}{3}
$$

$$
\iint_D f(x,y) \, dA = \int_a^b \int_c^d f(x,y) \, dy \, dx
$$

$$
\oint_C \vec{F} \cdot d\vec{r}
$$

---

## 7. Summations and Products

$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
$$

$$
\prod_{i=1}^{n} i = n!
$$

---

## 8. Matrices

$$
\begin{pmatrix}
a & b \\
c & d
\end{pmatrix}
\begin{bmatrix}
x \\
y
\end{bmatrix}
=
\begin{bmatrix}
ax + by \\
cx + dy
\end{bmatrix}
$$

$$
\det(A) = \begin{vmatrix}
a_{11} & a_{12} & a_{13} \\
a_{21} & a_{22} & a_{23} \\
a_{31} & a_{32} & a_{33}
\end{vmatrix}
$$

---

## 9. Physics Equations

**Maxwell's Equations:**

$$
\nabla \cdot \vec{E} = \frac{\rho}{\epsilon_0}
$$

$$
\nabla \cdot \vec{B} = 0
$$

$$
\nabla \times \vec{E} = -\frac{\partial \vec{B}}{\partial t}
$$

$$
\nabla \times \vec{B} = \mu_0 \vec{J} + \mu_0 \epsilon_0 \frac{\partial \vec{E}}{\partial t}
$$

**Schrödinger Equation:**

$$
i\hbar \frac{\partial}{\partial t} \Psi(\vec{r}, t) = \hat{H} \Psi(\vec{r}, t)
$$

---

## 10. Chemistry Notation

$$
\ce{H2O} \quad \ce{CO2} \quad \ce{NaCl}
$$

$$
\ce{2H2 + O2 -> 2H2O}
$$

*(Note: `\ce{}` requires mhchem extension)*

Alternative without mhchem:

$$
\text{2H}_2 + \text{O}_2 \rightarrow \text{2H}_2\text{O}
$$

---

## 11. Limits and Calculus

$$
\lim_{x \to 0} \frac{\sin x}{x} = 1
$$

$$
\frac{d}{dx}\left[ \int_0^x f(t) \, dt \right] = f(x)
$$

$$
f'(x) = \lim_{h \to 0} \frac{f(x+h) - f(x)}{h}
$$

---

## 12. Special Functions

$$
\Gamma(n) = (n-1)! = \int_0^{\infty} t^{n-1} e^{-t} \, dt
$$

$$
\zeta(s) = \sum_{n=1}^{\infty} \frac{1}{n^s}
$$

---

## 13. Logic and Set Theory

$$
\forall x \in \mathbb{R}, \exists y \in \mathbb{R} : x + y = 0
$$

$$
A \cup B = \{x : x \in A \lor x \in B\}
$$

$$
A \cap B = \{x : x \in A \land x \in B\}
$$

$$
A \subseteq B \iff (x \in A \Rightarrow x \in B)
$$

---

## 14. Brackets and Delimiters

$$
\left( \frac{a}{b} \right) \quad \left[ \frac{a}{b} \right] \quad \left\{ \frac{a}{b} \right\}
$$

$$
\left\langle \psi \middle| \phi \right\rangle
$$

$$
\left\lfloor x \right\rfloor \quad \left\lceil x \right\rceil
$$

---

## 15. Aligned Equations

$$
\begin{align}
(a + b)^2 &= a^2 + 2ab + b^2 \\
(a - b)^2 &= a^2 - 2ab + b^2 \\
(a + b)(a - b) &= a^2 - b^2
\end{align}
$$

---

## 16. Cases / Piecewise Functions

$$
|x| = \begin{cases}
x & \text{if } x \geq 0 \\
-x & \text{if } x < 0
\end{cases}
$$

$$
f(n) = \begin{cases}
n/2 & \text{if } n \text{ is even} \\
3n+1 & \text{if } n \text{ is odd}
\end{cases}
$$

---

## 17. Decorations

$$
\hat{x} \quad \bar{x} \quad \vec{x} \quad \tilde{x} \quad \dot{x} \quad \ddot{x}
$$

$$
\overline{AB} \quad \underline{text} \quad \overbrace{a+b+c}^{\text{sum}} \quad \underbrace{x+y}_{\text{terms}}
$$

---

## 18. Colors (if supported)

$$
\textcolor{red}{x^2} + \textcolor{blue}{y^2} = \textcolor{green}{z^2}
$$

---

*End of KaTeX test document*
