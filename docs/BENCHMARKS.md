```fish startup benchmark
time begin
    fish -i -c exit 0
end
```

timed using the fish shell `time` builtin (results are the median of five runs)

### Startup:

```fish catbow
time echo asdf | cb
________________________________________________________
Executed in    3.90 millis    fish           external
   usr time   51.21 millis  533.00 micros   50.67 millis
   sys time   31.24 millis   96.00 micros   31.14 millis
```
```fish lolcat
time echo asdf | lolcat
________________________________________________________
Executed in  108.07 millis    fish           external
   usr time   93.23 millis  629.00 micros   92.60 millis
   sys time   14.11 millis  201.00 micros   13.91 millis
```

**28x improvement in startup time**

sample text is generated with generate_text.py and piped to a file 

### 100 lines, 60 character width
`./generate_text.py --line-width 60 --num-lines 100`

```fish catbow
Executed in   36.47 millis
```
```fish lolcat
Executed in  206.32 millis
```
**6x faster**

### 10k lines, 60 character width
`./generate_text.py --line-width 60 --num-lines 10_000`

```fish catbow
Executed in    1.18 secs
```
```fish lolcat
Executed in    9.62 secs
```
**8x faster**

### 10k lines, 1000 character width
`./generate_text.py --line-width 1_000 --num-lines 10_000`

```fish catbow
Executed in   18.91 secs
```
```fish lolcat
Executed in  169.76 secs
```
**111x faster**

### 10k lines, 1 character width
`./generate_text.py --line-width 1 --num-lines 10_000`

```fish catbow
Executed in   84.68 millis
```
```fish lolcat
Executed in  485.44 millis
```
**6x faster**

### 1 line, 100k character width
`./generate_text.py --line-width 100_000 --num-lines 1`

```fish catbow
Executed in  257.46 millis
```
```fish lolcat
Executed in    1.68 secs 
```
**7x faster**
