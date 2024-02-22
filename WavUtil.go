package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// ConvertPCMToWAV converts PCM data to WAV format and writes it to the given file
// The sampleRate, bitsPerSample, and numChannels parameters are used to populate the WAV header
// The PCM data is expected to be in big-endian format
// The temporary file is expected to be closed by the caller
func ConvertPCMToWAV(tmpFile *os.File, pcmData []byte, sampleRate, bitsPerSample, numChannels uint32) error {
	header := NewWaveHeader(sampleRate, bitsPerSample, numChannels, uint32(len(pcmData)))

	// Writes the header
	if err := binary.Write(tmpFile, binary.LittleEndian, header); err != nil {
		return fmt.Errorf("failed to write WAV header: %w", err)
	}

	// Writes the PCM data
	if _, err := tmpFile.Write(convertBigEndianToLittleEndian(pcmData)); err != nil {
		return fmt.Errorf("failed to write PCM data: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to flush writes to temporary file: %w", err)
	}

	return nil
}

type WaveHeader struct {
	RIFFHeader    [4]byte // Contains "RIFF"
	FileSize      uint32  // File size in bytes
	WAVEHeader    [4]byte // Contains "WAVE"
	FmtHeader     [4]byte // Contains "fmt "
	FmtChunkSize  uint32  // Size of the following fmt header
	AudioFormat   uint16  // Audio format, 1 for PCM
	NumChannels   uint16  // Number of channels
	SampleRate    uint32  // Sample rate
	ByteRate      uint32  // Byte rate
	BlockAlign    uint16  // Block align
	BitsPerSample uint16  // Bits per sample
	DataHeader    [4]byte // Contains "data"
	DataChunkSize uint32  // Size of the following data
}

func convertBigEndianToLittleEndian(data []byte) []byte {
	if len(data)%2 != 0 {
		return nil
	}
	converted := make([]byte, len(data))
	for i := 0; i < len(data); i += 2 {
		converted[i], converted[i+1] = data[i+1], data[i]
	}
	return converted
}
func NewWaveHeader(sampleRate, bitsPerSample, numChannels, dataSize uint32) WaveHeader {
	return WaveHeader{
		RIFFHeader:    [4]byte{'R', 'I', 'F', 'F'},
		FileSize:      36 + dataSize,
		WAVEHeader:    [4]byte{'W', 'A', 'V', 'E'},
		FmtHeader:     [4]byte{'f', 'm', 't', ' '},
		FmtChunkSize:  16,
		AudioFormat:   1,
		NumChannels:   uint16(numChannels),
		SampleRate:    sampleRate,
		ByteRate:      sampleRate * numChannels * (bitsPerSample / 8),
		BlockAlign:    uint16(numChannels * (bitsPerSample / 8)),
		BitsPerSample: uint16(bitsPerSample),
		DataHeader:    [4]byte{'d', 'a', 't', 'a'},
		DataChunkSize: dataSize,
	}
}
