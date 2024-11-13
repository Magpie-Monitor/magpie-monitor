import React, {useEffect, useState} from 'react';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent';
import ActionButton, { ActionButtonColor } from 'components/ActionButton/ActionButton.tsx';
import './CustomPromptPopup.scss';

interface CustomPromptPopupProps {
    initialValue: string;
    isDisplayed: boolean;
    onSave: (value: string) => void;
    onClose: () => void;
}

const CustomPromptPopup: React.FC<CustomPromptPopupProps> =
    ({ initialValue, isDisplayed, onSave, onClose }) => {
        const [customPrompt, setCustomPrompt] = useState(initialValue);

        useEffect(() => {
            setCustomPrompt(initialValue);
        }, [initialValue]);

        const handleSave = () => {
            onSave(customPrompt);
        };

    return (
        <OverlayComponent isDisplayed={isDisplayed} onClose={onClose}>
            <div className="custom-prompt-popup">
                <h3 className="custom-prompt-popup__title">Edit Custom Prompt</h3>
                <textarea
                    className="custom-prompt-popup__textarea"
                    value={customPrompt}
                    onChange={(e) => setCustomPrompt(e.target.value)}
                    placeholder="Enter your custom prompt here..."
                />
                <div className="custom-prompt-popup__actions">
                    <ActionButton onClick={handleSave}
                                  description="Save" color={ActionButtonColor.GREEN}/>
                    <ActionButton onClick={onClose}
                                  description="Cancel" color={ActionButtonColor.RED}/>
                </div>
            </div>
        </OverlayComponent>
    );
    };

export default CustomPromptPopup;
