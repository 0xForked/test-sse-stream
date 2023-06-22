package domain

const SequenceKey = "sequence"

const TTSChannelName = "TTS"

const PyScript = `
from gtts import gTTS
from playsound import playsound
import os

text = "%s"
tts = gTTS(text, lang='id')

output_file = "output.mp3"
tts.save(output_file)
playsound(output_file)
os.remove(output_file)
`
