import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import { useState } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import { AccuracyLevel } from 'api/managment-service';

const AccuracySection = () => {
    const [accuracy, setAccuracy] = useState<AccuracyLevel>('HIGH');
    return (
        <SectionComponent icon={<SVGIcon iconName='precision-icon' />} title={'Accuracy'}>
            <div className="precision-section__input-group">
                <TagButton
                    listItems={['HIGH', 'MEDIUM', 'LOW']}
                    chosenItem={accuracy}
                    onSelect={setAccuracy}
                />
            </div>
        </SectionComponent>
    );
};
export default AccuracySection;