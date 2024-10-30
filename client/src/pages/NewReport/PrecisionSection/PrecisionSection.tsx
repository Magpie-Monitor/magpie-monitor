import './PrecisionSection.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent.tsx';
import TagButton from 'components/TagButton/TagButton.tsx';
import { useState } from 'react';

const PrecisionSection = () => {
    const [state, setState] = useState('enabled');
    const [precision, setPrecision] = useState('high');

    return (
        <SectionComponent icon={'precision-icon'} title={'Precision'}>
            <div className="precision-section__input-group">
                <p>State</p>
                <TagButton
                    listItems={['enabled', 'disabled']}
                    chosenItem={state}
                    onSelect={setState}
                />
                <p>Precision</p>
                <TagButton
                    listItems={['high', 'low']}
                    chosenItem={precision}
                    onSelect={setPrecision}
                />
            </div>
        </SectionComponent>
    );
};

export default PrecisionSection;
