#include<stdio.h>
#include<stdlib.h>
#include "_cgo_export.h"

#define BMP_HEADER_SIZE 54
int restHeaderSize;
int c_reader(void* callback, char *src) {
	FILE* f = fopen(src, "rb");

	if (f == NULL) {
		fclose(f);
		return 1;
	}

	unsigned char info[BMP_HEADER_SIZE];
	fread(info, sizeof(unsigned char), BMP_HEADER_SIZE, f); // read the BMP_HEADER_SIZE-byte header

	// extract image height and width from header
	int width = *(int*) &info[18];
	int height = *(int*) &info[22];
	int infoSize = *(int*) &info[2];
	int row_padded = ((width * 3) + 3) & (~3);

	
	restHeaderSize = (infoSize - row_padded * height) - BMP_HEADER_SIZE;
	if (restHeaderSize > 0) {
		restHeaderSize = 68;
		unsigned char restHeader[restHeaderSize];
		fread(restHeader, sizeof(unsigned char), restHeaderSize, f);
	}

	int size = row_padded * height;

	getPixelInfo(callback, width, height);

	long sum = 0;
	unsigned char dataRow[row_padded];
	for (int yIdx = 0; yIdx < height; yIdx++) {
		fread(dataRow, sizeof(unsigned char), row_padded, f);
		for (int xIdx = 0; xIdx < width; xIdx++) {
			int xIdxColor = xIdx * 3;
			handelPixelData(callback, dataRow[xIdxColor],
					dataRow[xIdxColor + 1], dataRow[xIdxColor + 2]);
		}

	}

	fclose(f);
	return 0;
}

int c_writeDataToImageFile(int *data, char *fileName, int width, int height) {
	FILE *f;
	
	unsigned char bmpfileheader[14] = { 'B', 'M', 0, 0, 0, 0, 0, 0, 0, 0,
	BMP_HEADER_SIZE + restHeaderSize, 0, 0, 0 };
	unsigned char bmpinfoheader[40] = { 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
			0, 24, 0 };
	unsigned char bmppad[3] = { 0, 0, 0 };

	unsigned char *img = (unsigned char *) malloc(
			(3 * width * height) * sizeof(unsigned char));

	for (int j = 0; j < height * width * 3; ++j) {
		img[j] = (unsigned char) data[j];
	}

	int fileSize = width * height * 3 + BMP_HEADER_SIZE + restHeaderSize;

	bmpfileheader[2] = (unsigned char) (fileSize);
	bmpfileheader[3] = (unsigned char) (fileSize >> 8);
	bmpfileheader[4] = (unsigned char) (fileSize >> 16);
	bmpfileheader[5] = (unsigned char) (fileSize >> 24);

	bmpinfoheader[4] = (unsigned char) (width);
	bmpinfoheader[5] = (unsigned char) (width >> 8);
	bmpinfoheader[6] = (unsigned char) (width >> 16);
	bmpinfoheader[7] = (unsigned char) (width >> 24);
	bmpinfoheader[8] = (unsigned char) (height);
	bmpinfoheader[9] = (unsigned char) (height >> 8);
	bmpinfoheader[10] = (unsigned char) (height >> 16);
	bmpinfoheader[11] = (unsigned char) (height >> 24);

	f = fopen(fileName, "wb");

	if (f == NULL) {
		fclose(f);
		return 1;
	}

	fwrite(bmpfileheader, 1, 14, f);
	fwrite(bmpinfoheader, 1, 40, f);
	printf("Hallo ##### %d", restHeaderSize);
	if(restHeaderSize > 0) {
		unsigned char restHeader[restHeaderSize];
		for(int i = 0; i < restHeaderSize; i++) {
			restHeader[i] = 0;
		}
		
		fwrite(restHeader, 1, restHeaderSize, f);
	}

	for (int i = 0; i < height; i++) {
		fwrite(img + width * i * 3, 3, width, f);
	}

	free(img);
	fclose(f);

	return 0;
}
