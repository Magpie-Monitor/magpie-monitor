import './TagButton.scss';
import React, {useEffect, useRef, useReducer, ReactElement} from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

interface TagButtonProps<T extends ReactElement | string | number> {
  listItems: T[];
  chosenItem: T;
  onSelect: (item: T) => void;
}

type Action<T extends ReactElement | string | number> =
    | { type: 'TOGGLE' }
    | { type: 'CLOSE' }
    | { type: 'SELECT_OPTION'; payload: T };

interface State<T extends ReactElement | string | number> {
  isOpen: boolean;
  selectedOption: T;
}

function reducer<T extends ReactElement | string | number>
    (state: State<T>, action: Action<T>): State<T> {
  switch (action.type) {
    case 'TOGGLE':
      return { ...state, isOpen: !state.isOpen };
    case 'CLOSE':
      return { ...state, isOpen: false };
    case 'SELECT_OPTION':
      return { isOpen: false, selectedOption: action.payload };
    default:
      return state;
  }
}

const TagButton = <T extends ReactElement | string | number>
    ({ listItems, chosenItem, onSelect }: TagButtonProps<T>): React.ReactNode => {
    const [{ isOpen, selectedOption }, dispatch] = useReducer(reducer, {
        isOpen: false,
        selectedOption: chosenItem,
    });

    const menuRef = useRef<HTMLUListElement>(null);
    const buttonRef = useRef<HTMLButtonElement>(null);

    const closeMenu = () => dispatch({type: 'CLOSE'});

    useEffect(() => {

        const handleClickOutside =
            (event: MouseEvent) => {
                if (
                    menuRef.current &&
                    !menuRef.current.contains(event.target as Node) &&
                    !buttonRef.current?.contains(event.target as Node)
                ) {
                    closeMenu();
                }
            };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Escape') closeMenu();
    };

    const handleSelect = (item: T) => {
        dispatch({type: 'SELECT_OPTION', payload: item});
        onSelect(item);
        buttonRef.current?.focus();
    };

    return (
        <div className="tag-button">
            <button
                className="tag-button__toggle"
                ref={buttonRef}
                onClick={() => dispatch({type: 'TOGGLE'})}
                onKeyDown={handleKeyDown}
                aria-haspopup="listbox"
                aria-expanded={isOpen}
            >
                <span className="tag-button__toggle__description">
                    {selectedOption}
                </span>
                <SVGIcon iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'}/>
            </button>
            {isOpen && (
                <ul
                    className="tag-button__menu"
                    ref={menuRef}
                    role="listbox"
                    tabIndex={-1}
                >
                    {listItems.map((item, index) => (
                        <li
                            id={`tag-button-option-${index}`}
                            key={index}
                            className="tag-button__menu__element"
                            onClick={() => handleSelect(item)}
                            role="option"
                            tabIndex={0}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter') handleSelect(item);
                            }}
                        >
                            {item}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default TagButton;
