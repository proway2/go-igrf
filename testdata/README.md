Each file contains IGRF values at yearly intervals at one location. Location is set on the first line of the file.

Sample file:

```
Lat  59.900 geodetic    Long   39.900     0.000 km                     
   DATE       D   SV      I  SV      H    SV       X    SV       Y    SV       Z    SV      F    SV
 1900.5    5.55     4  71.01   0  16646   -10   16568   -12    1610    18   48357   -16  51142   -18
 1901.5    5.62     4  71.01   0  16636   -10   16556   -12    1629    18   48341   -16  51124   -18
 1902.5    5.68     4  71.02   0  16626   -10   16544   -12    1647    18   48326   -16  51106   -18

 ...

 2022.5   14.55     9  74.34   3  14578   -19   14111   -28    3662    32   51988    80  53993    72
 2023.5   14.70     9  74.38   3  14559   -19   14083   -28    3694    32   52068    80  54065    72
 2024.5   14.85     9  74.42   3  14540   -19   14055   -28    3726    32   52148    80  54137    72
```
 
Where:
- D is declination in degrees (+ve east)
- I is inclination in degrees (+ve down)
- H is horizontal intensity in nT
- X is north component in nT
- Y is east component in nT
- Z is vertical component in nT (+ve down)
- F is total intensity in nT
- SV is secular variation (annual rate of change)
- Units for SV: minutes/yr (D & I); nT/yr (H,X,Y,Z & F)
