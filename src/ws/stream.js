const { SpeechClient } = require('@google-cloud/speech');

const stream = ( socket ) => {
    const speech = new SpeechClient();
    const fs = require('fs');

    socket.on( 'subscribe', ( data ) => {
        //subscribe/join a room
        socket.join( data.room );
        socket.join( data.socketId );

        //Inform other members in the room of new user's arrival
        if ( socket.adapter.rooms.has(data.room) === true ) {
            socket.to( data.room ).emit( 'new user', { socketId: data.socketId } );
        }
    } );


    socket.on( 'newUserStart', ( data ) => {
        socket.to( data.to ).emit( 'newUserStart', { sender: data.sender } );
    } );


    socket.on( 'sdp', ( data ) => {
        socket.to( data.to ).emit( 'sdp', { description: data.description, sender: data.sender } );
    } );


    socket.on( 'ice candidates', ( data ) => {
        socket.to( data.to ).emit( 'ice candidates', { candidate: data.candidate, sender: data.sender } );
    } );


    socket.on( 'chat', ( data ) => {
        socket.to( data.room ).emit( 'chat', { sender: data.sender, msg: data.msg } );
    } );

    socket.on('audio', async (data) => {
        try {
            const audioBuffer = Buffer.from(data.audio, 'base64');

            const audio = {
                content: audioBuffer,
            };

            const config = {
                encoding: 'LINEAR16',
                sampleRateHertz: 16000,
                languageCode: 'en-US',
            };

            const request = {
                audio: audio,
                config: config,
            };

            // Performs speech recognition on the audio file
            const [response] = await speech.recognize(request);
            const transcription = response.results
                .map((result) => result.alternatives[0].transcript)
                .join('\n');

            // Emit the transcription to the client
            socket.emit('transcription', { transcription });
        } catch (error) {
            console.error('Error transcribing audio:', error);
            socket.emit('transcriptionError', { error: error.message });
        }
    });
};

module.exports = stream;
