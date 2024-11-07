import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import React, { useState } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import { AccuracyLevel } from 'api/managment-service';

interface AccuracySectionProps {
    setParentAccuracy: (accuracy: AccuracyLevel) => void;
}

const AccuracySection: React.FC<AccuracySectionProps> = ({ setParentAccuracy }) => {
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');

    const handleAccuracyChange = (newAccuracy: AccuracyLevel) => {
        setAccuracy(newAccuracy);
        setParentAccuracy(newAccuracy);
    };

    return (
        <SectionComponent icon={<SVGIcon iconName='precision-icon' />} title={'Default accuracy'}>
            <div className="precision-section__input-group">
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={accuracy}
                    onSelect={handleAccuracyChange}
                />
            </div>
        </SectionComponent>
    );
};

export default AccuracySection;
