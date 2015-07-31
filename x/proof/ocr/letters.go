package main

var (
	inputs [][]float64
)

func init() {
	inputs = make([][]float64, 26)

	// Letter A
	inputs[0] = make([]float64, 35)
	inputs[0][1] = 1.000000
	inputs[0][5] = 1.000000
	inputs[0][6] = 1.000000
	inputs[0][10] = 1.000000
	inputs[0][11] = 1.000000
	inputs[0][15] = 1.000000
	inputs[0][16] = 1.000000
	inputs[0][17] = 1.000000
	inputs[0][18] = 1.000000
	inputs[0][19] = 1.000000
	inputs[0][20] = 1.000000
	inputs[0][21] = 1.000000
	inputs[0][25] = 1.000000
	inputs[0][26] = 1.000000
	inputs[0][30] = 1.000000
	inputs[0][32] = 1.000000
	inputs[0][33] = 1.000000
	inputs[0][34] = 1.000000

	// Letter B
	inputs[1] = make([]float64, 35)
	inputs[1][1] = 1.000000
	inputs[1][2] = 1.000000
	inputs[1][3] = 1.000000
	inputs[1][4] = 1.000000
	inputs[1][7] = 1.000000
	inputs[1][10] = 1.000000
	inputs[1][12] = 1.000000
	inputs[1][15] = 1.000000
	inputs[1][17] = 1.000000
	inputs[1][18] = 1.000000
	inputs[1][19] = 1.000000
	inputs[1][22] = 1.000000
	inputs[1][25] = 1.000000
	inputs[1][27] = 1.000000
	inputs[1][30] = 1.000000
	inputs[1][31] = 1.000000
	inputs[1][32] = 1.000000
	inputs[1][33] = 1.000000
	inputs[1][34] = 1.000000

	// Letter C
	inputs[2] = make([]float64, 35)
	inputs[2][2] = 1.000000
	inputs[2][3] = 1.000000
	inputs[2][4] = 1.000000
	inputs[2][5] = 1.000000
	inputs[2][6] = 1.000000
	inputs[2][11] = 1.000000
	inputs[2][16] = 1.000000
	inputs[2][21] = 1.000000
	inputs[2][26] = 1.000000
	inputs[2][32] = 1.000000
	inputs[2][33] = 1.000000
	inputs[2][34] = 1.000000

	// Letter D
	inputs[3] = make([]float64, 35)
	inputs[3][1] = 1.000000
	inputs[3][2] = 1.000000
	inputs[3][3] = 1.000000
	inputs[3][4] = 1.000000
	inputs[3][7] = 1.000000
	inputs[3][10] = 1.000000
	inputs[3][12] = 1.000000
	inputs[3][15] = 1.000000
	inputs[3][17] = 1.000000
	inputs[3][20] = 1.000000
	inputs[3][22] = 1.000000
	inputs[3][25] = 1.000000
	inputs[3][27] = 1.000000
	inputs[3][30] = 1.000000
	inputs[3][31] = 1.000000
	inputs[3][32] = 1.000000
	inputs[3][33] = 1.000000
	inputs[3][34] = 1.000000

	// Letter E
	inputs[4] = make([]float64, 35)
	inputs[4][1] = 1.000000
	inputs[4][2] = 1.000000
	inputs[4][3] = 1.000000
	inputs[4][4] = 1.000000
	inputs[4][5] = 1.000000
	inputs[4][6] = 1.000000
	inputs[4][11] = 1.000000
	inputs[4][16] = 1.000000
	inputs[4][17] = 1.000000
	inputs[4][18] = 1.000000
	inputs[4][19] = 1.000000
	inputs[4][21] = 1.000000
	inputs[4][26] = 1.000000
	inputs[4][31] = 1.000000
	inputs[4][32] = 1.000000
	inputs[4][33] = 1.000000
	inputs[4][34] = 1.000000

	// Letter F
	inputs[5] = make([]float64, 35)
	inputs[5][1] = 1.000000
	inputs[5][6] = 1.000000
	inputs[5][11] = 1.000000
	inputs[5][16] = 1.000000
	inputs[5][17] = 1.000000
	inputs[5][18] = 1.000000
	inputs[5][19] = 1.000000
	inputs[5][21] = 1.000000
	inputs[5][26] = 1.000000
	inputs[5][31] = 1.000000
	inputs[5][32] = 1.000000
	inputs[5][33] = 1.000000
	inputs[5][34] = 1.000000

	// Letter G
	inputs[6] = make([]float64, 35)
	inputs[6][2] = 1.000000
	inputs[6][3] = 1.000000
	inputs[6][4] = 1.000000
	inputs[6][5] = 1.000000
	inputs[6][6] = 1.000000
	inputs[6][10] = 1.000000
	inputs[6][11] = 1.000000
	inputs[6][15] = 1.000000
	inputs[6][16] = 1.000000
	inputs[6][19] = 1.000000
	inputs[6][20] = 1.000000
	inputs[6][21] = 1.000000
	inputs[6][26] = 1.000000
	inputs[6][32] = 1.000000
	inputs[6][33] = 1.000000
	inputs[6][34] = 1.000000

	// Letter H
	inputs[7] = make([]float64, 35)
	inputs[7][1] = 1.000000
	inputs[7][5] = 1.000000
	inputs[7][6] = 1.000000
	inputs[7][10] = 1.000000
	inputs[7][11] = 1.000000
	inputs[7][15] = 1.000000
	inputs[7][16] = 1.000000
	inputs[7][17] = 1.000000
	inputs[7][18] = 1.000000
	inputs[7][19] = 1.000000
	inputs[7][20] = 1.000000
	inputs[7][21] = 1.000000
	inputs[7][25] = 1.000000
	inputs[7][26] = 1.000000
	inputs[7][30] = 1.000000
	inputs[7][31] = 1.000000

	// Letter I
	inputs[8] = make([]float64, 35)
	inputs[8][3] = 1.000000
	inputs[8][8] = 1.000000
	inputs[8][13] = 1.000000
	inputs[8][18] = 1.000000
	inputs[8][23] = 1.000000
	inputs[8][28] = 1.000000
	inputs[8][33] = 1.000000

	// Letter J
	inputs[9] = make([]float64, 35)
	inputs[9][1] = 1.000000
	inputs[9][2] = 1.000000
	inputs[9][3] = 1.000000
	inputs[9][4] = 1.000000
	inputs[9][6] = 1.000000
	inputs[9][10] = 1.000000
	inputs[9][15] = 1.000000
	inputs[9][20] = 1.000000
	inputs[9][25] = 1.000000
	inputs[9][30] = 1.000000
	inputs[9][33] = 1.000000
	inputs[9][34] = 1.000000

	// Letter K
	inputs[10] = make([]float64, 35)
	inputs[10][1] = 1.000000
	inputs[10][5] = 1.000000
	inputs[10][6] = 1.000000
	inputs[10][9] = 1.000000
	inputs[10][10] = 1.000000
	inputs[10][11] = 1.000000
	inputs[10][13] = 1.000000
	inputs[10][14] = 1.000000
	inputs[10][16] = 1.000000
	inputs[10][17] = 1.000000
	inputs[10][18] = 1.000000
	inputs[10][21] = 1.000000
	inputs[10][23] = 1.000000
	inputs[10][24] = 1.000000
	inputs[10][26] = 1.000000
	inputs[10][29] = 1.000000
	inputs[10][30] = 1.000000
	inputs[10][31] = 1.000000

	// Letter L
	inputs[11] = make([]float64, 35)
	inputs[11][1] = 1.000000
	inputs[11][2] = 1.000000
	inputs[11][3] = 1.000000
	inputs[11][4] = 1.000000
	inputs[11][5] = 1.000000
	inputs[11][6] = 1.000000
	inputs[11][11] = 1.000000
	inputs[11][16] = 1.000000
	inputs[11][21] = 1.000000
	inputs[11][26] = 1.000000
	inputs[11][31] = 1.000000

	// Letter M
	inputs[12] = make([]float64, 35)
	inputs[12][1] = 1.000000
	inputs[12][5] = 1.000000
	inputs[12][6] = 1.000000
	inputs[12][10] = 1.000000
	inputs[12][11] = 1.000000
	inputs[12][15] = 1.000000
	inputs[12][16] = 1.000000
	inputs[12][18] = 1.000000
	inputs[12][20] = 1.000000
	inputs[12][21] = 1.000000
	inputs[12][22] = 1.000000
	inputs[12][23] = 1.000000
	inputs[12][24] = 1.000000
	inputs[12][25] = 1.000000
	inputs[12][26] = 1.000000
	inputs[12][27] = 1.000000
	inputs[12][29] = 1.000000
	inputs[12][30] = 1.000000
	inputs[12][31] = 1.000000

	// Letter N
	inputs[13] = make([]float64, 35)
	inputs[13][1] = 1.000000
	inputs[13][5] = 1.000000
	inputs[13][6] = 1.000000
	inputs[13][9] = 1.000000
	inputs[13][10] = 1.000000
	inputs[13][11] = 1.000000
	inputs[13][13] = 1.000000
	inputs[13][14] = 1.000000
	inputs[13][15] = 1.000000
	inputs[13][16] = 1.000000
	inputs[13][18] = 1.000000
	inputs[13][20] = 1.000000
	inputs[13][21] = 1.000000
	inputs[13][22] = 1.000000
	inputs[13][23] = 1.000000
	inputs[13][25] = 1.000000
	inputs[13][26] = 1.000000
	inputs[13][27] = 1.000000
	inputs[13][30] = 1.000000
	inputs[13][31] = 1.000000

	// Letter O
	inputs[14] = make([]float64, 35)
	inputs[14][1] = 1.000000
	inputs[14][2] = 1.000000
	inputs[14][3] = 1.000000
	inputs[14][4] = 1.000000
	inputs[14][5] = 1.000000
	inputs[14][6] = 1.000000
	inputs[14][10] = 1.000000
	inputs[14][11] = 1.000000
	inputs[14][15] = 1.000000
	inputs[14][16] = 1.000000
	inputs[14][20] = 1.000000
	inputs[14][21] = 1.000000
	inputs[14][25] = 1.000000
	inputs[14][26] = 1.000000
	inputs[14][30] = 1.000000
	inputs[14][31] = 1.000000
	inputs[14][32] = 1.000000
	inputs[14][33] = 1.000000
	inputs[14][34] = 1.000000

	// Letter P
	inputs[15] = make([]float64, 35)
	inputs[15][2] = 1.000000
	inputs[15][7] = 1.000000
	inputs[15][12] = 1.000000
	inputs[15][17] = 1.000000
	inputs[15][18] = 1.000000
	inputs[15][19] = 1.000000
	inputs[15][22] = 1.000000
	inputs[15][25] = 1.000000
	inputs[15][27] = 1.000000
	inputs[15][30] = 1.000000
	inputs[15][31] = 1.000000
	inputs[15][32] = 1.000000
	inputs[15][33] = 1.000000
	inputs[15][34] = 1.000000

	// Letter Q
	inputs[16] = make([]float64, 35)
	inputs[16][2] = 1.000000
	inputs[16][3] = 1.000000
	inputs[16][5] = 1.000000
	inputs[16][6] = 1.000000
	inputs[16][9] = 1.000000
	inputs[16][11] = 1.000000
	inputs[16][13] = 1.000000
	inputs[16][15] = 1.000000
	inputs[16][16] = 1.000000
	inputs[16][20] = 1.000000
	inputs[16][21] = 1.000000
	inputs[16][25] = 1.000000
	inputs[16][26] = 1.000000
	inputs[16][30] = 1.000000
	inputs[16][32] = 1.000000
	inputs[16][33] = 1.000000
	inputs[16][34] = 1.000000

	// Letter R
	inputs[17] = make([]float64, 35)
	inputs[17][1] = 1.000000
	inputs[17][5] = 1.000000
	inputs[17][6] = 1.000000
	inputs[17][9] = 1.000000
	inputs[17][11] = 1.000000
	inputs[17][13] = 1.000000
	inputs[17][16] = 1.000000
	inputs[17][17] = 1.000000
	inputs[17][18] = 1.000000
	inputs[17][19] = 1.000000
	inputs[17][21] = 1.000000
	inputs[17][25] = 1.000000
	inputs[17][26] = 1.000000
	inputs[17][30] = 1.000000
	inputs[17][31] = 1.000000
	inputs[17][32] = 1.000000
	inputs[17][33] = 1.000000
	inputs[17][34] = 1.000000

	// Letter S
	inputs[18] = make([]float64, 35)
	inputs[18][1] = 1.000000
	inputs[18][2] = 1.000000
	inputs[18][3] = 1.000000
	inputs[18][4] = 1.000000
	inputs[18][10] = 1.000000
	inputs[18][15] = 1.000000
	inputs[18][17] = 1.000000
	inputs[18][18] = 1.000000
	inputs[18][19] = 1.000000
	inputs[18][21] = 1.000000
	inputs[18][26] = 1.000000
	inputs[18][32] = 1.000000
	inputs[18][33] = 1.000000
	inputs[18][34] = 1.000000

	// Letter T
	inputs[19] = make([]float64, 35)
	inputs[19][3] = 1.000000
	inputs[19][8] = 1.000000
	inputs[19][13] = 1.000000
	inputs[19][18] = 1.000000
	inputs[19][23] = 1.000000
	inputs[19][28] = 1.000000
	inputs[19][31] = 1.000000
	inputs[19][32] = 1.000000
	inputs[19][33] = 1.000000
	inputs[19][34] = 1.000000

	// Letter U
	inputs[20] = make([]float64, 35)
	inputs[20][2] = 1.000000
	inputs[20][3] = 1.000000
	inputs[20][4] = 1.000000
	inputs[20][6] = 1.000000
	inputs[20][10] = 1.000000
	inputs[20][11] = 1.000000
	inputs[20][15] = 1.000000
	inputs[20][16] = 1.000000
	inputs[20][20] = 1.000000
	inputs[20][21] = 1.000000
	inputs[20][25] = 1.000000
	inputs[20][26] = 1.000000
	inputs[20][30] = 1.000000
	inputs[20][31] = 1.000000

	// Letter V
	inputs[21] = make([]float64, 35)
	inputs[21][1] = 1.000000
	inputs[21][2] = 1.000000
	inputs[21][6] = 1.000000
	inputs[21][7] = 1.000000
	inputs[21][8] = 1.000000
	inputs[21][11] = 1.000000
	inputs[21][13] = 1.000000
	inputs[21][16] = 1.000000
	inputs[21][18] = 1.000000
	inputs[21][19] = 1.000000
	inputs[21][21] = 1.000000
	inputs[21][24] = 1.000000
	inputs[21][26] = 1.000000
	inputs[21][29] = 1.000000
	inputs[21][30] = 1.000000
	inputs[21][31] = 1.000000

	// Letter W
	inputs[22] = make([]float64, 35)
	inputs[22][2] = 1.000000
	inputs[22][4] = 1.000000
	inputs[22][6] = 1.000000
	inputs[22][8] = 1.000000
	inputs[22][10] = 1.000000
	inputs[22][11] = 1.000000
	inputs[22][13] = 1.000000
	inputs[22][15] = 1.000000
	inputs[22][16] = 1.000000
	inputs[22][18] = 1.000000
	inputs[22][20] = 1.000000
	inputs[22][21] = 1.000000
	inputs[22][25] = 1.000000
	inputs[22][26] = 1.000000
	inputs[22][30] = 1.000000
	inputs[22][31] = 1.000000

	// Letter X
	inputs[23] = make([]float64, 35)
	inputs[23][1] = 1.000000
	inputs[23][5] = 1.000000
	inputs[23][6] = 1.000000
	inputs[23][7] = 1.000000
	inputs[23][9] = 1.000000
	inputs[23][10] = 1.000000
	inputs[23][12] = 1.000000
	inputs[23][13] = 1.000000
	inputs[23][14] = 1.000000
	inputs[23][18] = 1.000000
	inputs[23][22] = 1.000000
	inputs[23][23] = 1.000000
	inputs[23][24] = 1.000000
	inputs[23][26] = 1.000000
	inputs[23][27] = 1.000000
	inputs[23][29] = 1.000000
	inputs[23][30] = 1.000000
	inputs[23][31] = 1.000000

	// Letter Y
	inputs[24] = make([]float64, 35)
	inputs[24][3] = 1.000000
	inputs[24][8] = 1.000000
	inputs[24][13] = 1.000000
	inputs[24][17] = 1.000000
	inputs[24][18] = 1.000000
	inputs[24][19] = 1.000000
	inputs[24][22] = 1.000000
	inputs[24][24] = 1.000000
	inputs[24][26] = 1.000000
	inputs[24][27] = 1.000000
	inputs[24][29] = 1.000000
	inputs[24][30] = 1.000000
	inputs[24][31] = 1.000000

	// Letter Z
	inputs[25] = make([]float64, 35)
	inputs[25][1] = 1.000000
	inputs[25][2] = 1.000000
	inputs[25][3] = 1.000000
	inputs[25][4] = 1.000000
	inputs[25][5] = 1.000000
	inputs[25][7] = 1.000000
	inputs[25][12] = 1.000000
	inputs[25][13] = 1.000000
	inputs[25][18] = 1.000000
	inputs[25][23] = 1.000000
	inputs[25][24] = 1.000000
	inputs[25][29] = 1.000000
	inputs[25][31] = 1.000000
	inputs[25][32] = 1.000000
	inputs[25][33] = 1.000000
	inputs[25][34] = 1.000000

}
