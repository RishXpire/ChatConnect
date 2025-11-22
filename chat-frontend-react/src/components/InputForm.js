import React, {useState} from 'react';

function InputForm({ onSendMessage }) {

    const [text, setText] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();

        const trimmedText = text.trim();
        if(trimmedText) {
            onSendMessage(trimmedText);

            setText('');
        }
    };

    return (
        <form className = "input-form" onSubmit = {handleSubmit}>
            <input 
                type = "text"
                placeholder = "Type a message..."
                value = {text}
                onChange = { (e) => setText(e.target.value)}
            />
            <button type = "submit"> Send </button>
        </form>
    );
}
 
export default InputForm;
