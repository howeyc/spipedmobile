package com.howeyc.spiped;

import android.app.Activity;
import android.os.Bundle;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.webkit.WebChromeClient;

import go.Go;
import go.spipedmobile.Spipedmobile;


public class MyActivity extends Activity {

    WebView browser;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.web_view);

        //Magic
        Go.init(getApplicationContext());

        Spipedmobile.Start();

        browser = (WebView)findViewById(R.id.browser);
        browser.getSettings().setJavaScriptEnabled(true);
        browser.setWebViewClient(new WebViewClient());
        browser.setWebChromeClient(new WebChromeClient());
        browser.loadUrl("http://localhost:56056");
    }
}
