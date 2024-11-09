import './TagButton.scss';
import React, { useEffect, useRef, useReducer, useCallback } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

interface TagButtonProps<T> {
  listItems: T[];
  chosenItem: T;
  onSelect: (item: T) => void;
}

type Action<T> =
    | { type: 'TOGGLE' }
    | { type: 'CLOSE' }
    | { type: 'SELECT_OPTION'; payload: T };

interface State<T> {
  isOpen: boolean;
  selectedOption: T;
}

function reducer<T>(state: State<T>, action: Action<T>): State<T> {
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

const TagButton = <T,>({ listItems, chosenItem, onSelect }: TagButtonProps<T>): JSX.Element => {
    const [{ isOpen, selectedOption }, dispatch] = useReducer(reducer, {
        isOpen: false,
        selectedOption: chosenItem,
    });

    const menuRef = useRef<HTMLUListElement>(null);
    const buttonRef = useRef<HTMLButtonElement>(null);

    const closeMenu = useCallback(() => dispatch({ type: 'CLOSE' }), []);

    const handleClickOutside = useCallback(
        (event: MouseEvent) => {
            if (
                menuRef.current &&
                !menuRef.current.contains(event.target as Node) &&
                !buttonRef.current?.contains(event.target as Node)
            ) {
                closeMenu();
            }
        },
        [closeMenu]
    );

    useEffect(() => {
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [handleClickOutside]);

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Escape') closeMenu();
    };

    const handleSelect = (item: T) => {
        dispatch({ type: 'SELECT_OPTION', payload: item });
        onSelect(item);
        buttonRef.current?.focus();
    };

    return (
        <div className="tag-button">
            <button
                className="tag-button__toggle"
                ref={buttonRef}
                onClick={() => dispatch({ type: 'TOGGLE' })}
                onKeyDown={handleKeyDown}
                aria-haspopup="listbox"
                aria-expanded={isOpen}
            >
                <span className="tag-button__toggle__description">
                    {selectedOption as string}
                </span>
                <SVGIcon iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'} />
            </button>
            {isOpen && (
                <ul
                    className="tag-button__menu"
                    ref={menuRef}
                    role="listbox"
                    aria-activedescendant={`tag-button-option-${listItems.indexOf(
                        selectedOption as T
                    )}`}
                    tabIndex={-1}
                >
                    {listItems.map((item, index) => (
                        <li
                            id={`tag-button-option-${index}`}
                            key={index}
                            className="tag-button__menu__element"
                            onClick={() => handleSelect(item)}
                            role="option"
                            aria-selected={selectedOption === item}
                            tabIndex={0}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' || e.key === ' ') handleSelect(item);
                            }}
                        >
                            {item as string}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default TagButton;
