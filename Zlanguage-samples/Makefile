CXX ?= c++
CXXFLAGS ?= -std=c++17 -Wall -Wextra -pedantic

TARGET := makefile-smoke
SOURCE := test81.cpp

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SOURCE)
	$(CXX) $(CXXFLAGS) $(SOURCE) -o $(TARGET)

clean:
	$(RM) $(TARGET)
