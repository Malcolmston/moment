import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { DocsView } from '../../../src/components/DocsView';
import type { DocIndex } from 'go-ui';

// A minimal DocIndex the stubbed fetch returns for DocsApp's doc.json request.
const DOC_INDEX: DocIndex = {
  module: 'github.com/malcolmston/moment',
  packages: [
    {
      importPath: 'github.com/malcolmston/moment',
      name: 'moment',
      synopsis: 'Package moment is a standard-library-only moment.js-style date/time API over time.',
      doc: 'Package moment is a standard-library-only moment.js-style date/time API over time.',
      consts: [],
      vars: [],
      types: [
        {
          name: 'Moment',
          signature: 'type Moment struct{}',
          doc: 'Moment is an immutable wrapper around time.Time.',
          consts: [],
          vars: [],
          funcs: [],
          methods: [],
        },
      ],
      funcs: [{ name: 'ParseFormat', signature: 'func ParseFormat(value, format string) (Moment, error)', doc: 'ParseFormat parses value using a moment-style token format.' }],
    },
  ],
};

describe('DocsView', () => {
  beforeEach(() => {
    // DocsApp fetches doc.json; return the small index.
    global.fetch = vi.fn((input: RequestInfo | URL) => {
      if (String(input).includes('doc.json')) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve(DOC_INDEX) } as Response);
      }
      return new Promise<Response>(() => {});
    }) as unknown as typeof fetch;
  });

  it('renders the inline React API reference from the fetched doc.json', async () => {
    const { container } = render(<DocsView />);
    expect(container.querySelector('#view-docs')).not.toBeNull();
    expect(
      screen.getByRole('heading', { level: 2, name: /API documentation/ }),
    ).toBeInTheDocument();

    // DocsApp fetches asynchronously, then renders the package view + symbols.
    expect(await screen.findByRole('heading', { name: /package moment/ })).toBeInTheDocument();
    expect(container.querySelector('#sym-ParseFormat'), 'func ParseFormat symbol card').not.toBeNull();
    expect(container.querySelector('#sym-Moment'), 'type Moment symbol card').not.toBeNull();

    // The secondary link to the raw generated static HTML remains.
    expect(screen.getByRole('link', { name: /Open the raw generated HTML/ })).toHaveAttribute('href', './api/');
  });
});
