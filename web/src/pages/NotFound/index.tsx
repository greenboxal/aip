import React from 'react'

const css: React.CSSProperties = {
    alignItems: 'center',
    display: 'flex',
    height: '100%',
    justifyContent: 'center',
    width: '100%',
}

 const NotFound: React.FC = () => (
    <div data-cy="not-found-page" style={css}>
        <h1>Page Not Found</h1>
    </div>
)

export default NotFound